#include <iostream>
#include <fstream>
#include <unordered_map>
#include <dirent.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <string.h>
#include <unistd.h>
#include <algorithm>
#include <vector>
#include "Persist.h"
#include <dlfcn.h>
#include "../Undo-Manager/Single_undo.h"
#include "../Undo-Manager/Multiple_undo.h"
using namespace std;

ofstream fout;

unordered_map<string, int> id2ver;
vector<tuple<string, int, string>> update_list;

extern "C"
{
    void Create_File(const char *id)
    {

        struct stat st = {0};
        if (stat("src/Operations/DBs/", &st) == -1)
        {
            mkdir("src/Operations/DBs/", 0700);
        }
        string filename = "src/Operations/DBs/" + string(id) + ".txt";
        fout.open(filename, ios_base::app);
        // cout << "VII" << endl;
        cout << "Storage file created for Replica " << id << endl;
        // cout << "VIII" << endl;
        fout.close();
        // cout << "IX" << endl;
        id2ver[string(id)] = 0;
    }

    void Record(const char *id, const char *val)
    {
        // Write to Physical file
        string filename = "src/Operations/DBs/" + string(id) + ".txt";
        id2ver[string(id)]++;
        fout.open(filename, ios_base::app);
        fout << id << " " << id2ver[string(id)] << " " << val << endl;
        fout.close();

        // End writing to Physical file
        tuple<string, int, string> t(string(id), id2ver[id], string(val));
        update_list.push_back(t);
    }

    const char *process_ret(string ret_val)
    {
        string roll_element(ret_val.substr(ret_val.rfind(" ")));
        roll_element.erase(remove(roll_element.begin(), roll_element.end(), ' '), roll_element.end());
        return roll_element.c_str();
    }

    int *process_vals(vector<tuple<string, int, string>> ul, int prev_updates)
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

    const char *undo(const char *id, int opt_nums)
    {
        string filename = "src/Operations/DBs/" + string(id) + ".txt";
        string ret_val;
        // Undo based on data corruption detection
        if (opt_nums == -2)
        {
            if (db_data_corrupt_undo(3)) // 3 prev updates for example purpose, can be configured
            {
                ret_val = multiple_undo(filename.c_str(), 3);
                return process_ret(ret_val);
            }
            else
            {
                cout << "no undo required based on data corruption" << endl;
                return NULL;
            }
        }
        // Undo single update operations
        if (opt_nums == 1)
        {
            cout << "Undoing the last update operations" << endl;
            // ret_val = single_undo(filename.c_str(), opt_nums);
            // hack
            ret_val = get<2>(update_list[update_list.size() - 2]);
        }
        else
        {
            cout << "Undoing the last " << opt_nums << " update operations" << endl;
            // ret_val = multiple_undo(filename.c_str(), opt_nums);
            // hack
            ret_val = get<2>(update_list[update_list.size() - (opt_nums + 1)]);
        }

        return process_ret(ret_val);
    }

    bool db_data_corrupt_undo(int prev_updates)
    {
        void *handle = dlopen("src/Operations/Plug-in/sha256.so", RTLD_NOW);

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

        int *data = process_vals(update_list, prev_updates);
        return should_undo_stat(data, prev_updates);
    }
}