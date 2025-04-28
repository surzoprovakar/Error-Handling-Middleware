#include <iostream>
#include <fstream>
#include <unordered_map>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <algorithm>
#include <vector>
#include "Persist.h"
#include <dlfcn.h>
#include "../Undo-Manager/Single_undo.h"
#include "../Undo-Manager/Multiple_undo.h"
#include "../Lamport/Lamport_Clock.h"
using namespace std;

ofstream fout;

unordered_map<int, int> id2ver;
vector<tuple<int, int, string>> update_list;
LamportClock c_clock;
unordered_map<int, LamportClock> rid2LT;

extern "C"
{
    void Create_File(int id)
    {
        struct stat st = {0};
        if (stat("DBs/", &st) == -1)
        {
            mkdir("DBs/", 0700);
        }
        string filename = "DBs/" + to_string(id) + "_local.txt";
        fout.open(filename, ios_base::app);
        cout << "Storage file created for Replica " << id << endl;
        fout.close();

        id2ver[id] = 0;
    }

    void Record(int id, const char *val, const char *opt_name)
    {
        // Write to Physical file
        string filename = "DBs/" + to_string(id) + "_local.txt";
        id2ver[id]++;
        fout.open(filename, ios_base::app);
        fout << "R_ID:" << id
             << " Ver:" << id2ver[id]
             << " St:" << val
             << " LT:" << c_clock.local_event()
             << " Opt:" << opt_name
             << " Type:Local" << endl;

        // fout << id << " " << id2ver[id] << " " << val << endl;
        fout.close();
        // End writing to Physical file

        tuple<int, int, string> t(id, id2ver[id], val);
        update_list.push_back(t);
    }

    void Remote_Record(int id, int rid, const char *val, const char *opt_name)
    {
        string filename = "DBs/" + to_string(id) + "_local.txt";
        id2ver[id]++;
        fout.open(filename, ios_base::app);
        fout << "R_ID:" << id
             << " Ver:" << id2ver[id]
             << " St:" << val
             << " LT:" << c_clock.local_event()
             << " Opt:" << opt_name
             << " Type:{Remote, rid:" << rid
             << " local_LT:" << rid2LT[rid].local_event()
             << "}" << endl;
        // fout << "replica_" + to_string(rid) << " " << opt_name << " " << c_clock.local_event() << endl;

        // Record remote replica's local Lamport Timestamp
        // rid2LT[rid].local_event();
        // cout << "rid:" << rid << " LT:" << rid2LT[rid].get_time() << endl;
        fout.close();
    }

    const char *process_ret(string ret_val)
    {
        string roll_element(ret_val.substr(ret_val.rfind(" ")));
        cout << "going back to this value: " << roll_element << endl;
        roll_element.erase(remove(roll_element.begin(), roll_element.end(), ' '), roll_element.end());
        return roll_element.c_str();
    }

    int *process_vals(vector<tuple<int, int, string>> ul, int prev_updates)
    {
        int *lastvalues = new int[prev_updates];
        int index = 0;
        for (auto it = ul.rbegin(); it != ul.rend(); ++it)
        {
            string strVal = std::get<2>(*it);
            if (index < prev_updates)
            {
                lastvalues[index] = stoi(strVal);
                index++;
            }
            if (index >= prev_updates)
            {
                break;
            }
        }
        return lastvalues;
    }

    // Undo for Cundas
    /*const char *undo(int id, int opt_nums)
    {
        string filename = "DBs/" + to_string(id) + "_local.txt";
        string ret_val;

        // Undo single update operations
        if (opt_nums == 1)
        {
            cout << "Undoing the last update operations" << endl;
            ret_val = single_undo(filename.c_str(), opt_nums);
        }
        else
        {
            cout << "Undoing the last " << opt_nums << " update operations" << endl;
            ret_val = multiple_undo(filename.c_str(), opt_nums);
        }

        return process_ret(ret_val);
    }*/
    bool undo(int r_id, const char *undo_update, int r_LT)
    {
        if (static_cast<int>(rid2LT[r_id].get_time()) > r_LT)
        {
            return true;
        }
        return false;
    }
}
