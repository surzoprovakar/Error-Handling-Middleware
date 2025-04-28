#include <iostream>
#include <fstream>
#include <unordered_map>
#include <sys/types.h>
#include <sys/stat.h>
#include <unistd.h>
#include <algorithm>
#include <vector>
#include <tuple>
#include <sstream>
#include "Persist.h"
#include <dlfcn.h>
#include <regex>
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
        string filename = "DBs/" + to_string(id) + ".txt";
        fout.open(filename, ios_base::app);
        cout << "Storage file created for Replica " << id << endl;
        fout.close();

        id2ver[id] = 0;
    }

    void Record(int id, const char *val, const char *opt_name, int up_value)
    {
        // Write to Physical file
        string filename = "DBs/" + to_string(id) + ".txt";
        id2ver[id]++;
        fout.open(filename, ios_base::app);
        fout << "R_ID:" << id
             << " Ver:" << id2ver[id]
             << " St:" << val
             << " LT:" << c_clock.local_event()
             << " Opt:" << opt_name << "_" << up_value
             << " Type:Local" << endl;

        // fout << id << " " << id2ver[id] << " " << val << endl;
        fout.close();
        // End writing to Physical file

        tuple<int, int, string> t(id, id2ver[id], val);
        update_list.push_back(t);
    }

    void Remote_Record(int id, int rid, const char *opt_name, int up_value)
    {
        string filename = "DBs/" + to_string(id) + ".txt";
        fout.open(filename, ios_base::app);
        fout << "R_ID:" << id
             << " Opt:" << opt_name << "_" << up_value
             << " Type:{Remote, rid:" << rid
             << " local_LT:" << rid2LT[rid].local_event()
             << "}" << endl;
        // fout << "replica_" + to_string(rid) << " " << opt_name << " " << c_clock.local_event() << endl;

        // Record remote replica's local Lamport Timestamp
        // rid2LT[rid].local_event();
        // cout << "rid:" << rid << " LT:" << rid2LT[rid].get_time() << endl;
        fout.close();
    }

    /*
        string process_undo(string rev)
        {
            stringstream ss(rev);
            vector<string> tokens;
            string token;

            while (std::getline(ss, token, ' '))
            {
                tokens.push_back(token);
            }

            tokens.erase(tokens.begin(), tokens.begin() + 2);

            string res;
            for (const auto &t : tokens)
            {
                res += t + " ";
            }
            return res;
        }
        const char *process_ret(string ret_val)
        {
            // string roll_element(ret_val.substr(ret_val.rfind(" ")));
            string roll_element = process_undo(ret_val);
            cout << "going back to this value: " << roll_element << endl;
            // roll_element.erase(remove(roll_element.begin(), roll_element.end(), ' '), roll_element.end());
            return roll_element.c_str();
        }

        int *process_vals()
        {
            string last_str_val = string(last_val);

            // last_str_val.erase(remove(last_str_val.begin(), last_str_val.end(), ','), last_str_val.end());
            last_str_val.erase(remove(last_str_val.begin(), last_str_val.end(), ':'), last_str_val.end());
            last_str_val.erase(remove(last_str_val.begin(), last_str_val.end(), '{'), last_str_val.end());
            last_str_val.erase(remove(last_str_val.begin(), last_str_val.end(), '}'), last_str_val.end());
            last_str_val.erase(remove(last_str_val.begin(), last_str_val.end(), '\"'), last_str_val.end());

            stringstream ss(last_str_val);
            vector<int> intArray;
            int num;
            while (ss >> num)
            {
                intArray.push_back(num);
                if (ss.peek() == ',')
                {
                    ss.ignore();
                }
            }

            int *arr = new int[intArray.size()];

            for (size_t i = 0; i < intArray.size(); ++i)
            {
                arr[i] = intArray[i];
            }
            stat_size = intArray.size();
            return arr;
        }

        const char *undo(int id, int opt_nums)
        {

            string filename = "DBs/" + to_string(id) + ".txt";
            string ret_val;

            // Undo based on statistical configurations
            if (opt_nums == -1)
            {
                if (stat_undo(2))
                {
                    ret_val = multiple_undo(filename.c_str(), 2);
                    return process_ret(ret_val);
                }
                else
                {
                    cout << "no undo required" << endl;
                    return NULL;
                }
            }
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
        }

        bool stat_undo(int prev_updates)
        {
            void *handle = dlopen("Plug-in/statistics.so", RTLD_NOW);

            if (handle == NULL)
            {
                cout << "Error loading SO file: " << dlerror() << endl;
                return false;
            }

            typedef bool (*Should_undo)(int *, int);
            Should_undo should_undo_stat = (Should_undo)dlsym(handle, "should_undo");

            if (!should_undo_stat)
            {
                cout << "Error resolving symbol: " << dlerror() << endl;
                dlclose(handle);
                return false;
            }

            int *data = process_vals();
            cout << "stat size is " << stat_size << endl;
            return should_undo_stat(data, stat_size);
        }
        */

    bool undo(int r_id, const char *undo_update, int undo_val, int r_LT)
    {
        if (static_cast<int>(rid2LT[r_id].get_time()) > r_LT)
        {
            return true;
        }
        return false;
    }
}
