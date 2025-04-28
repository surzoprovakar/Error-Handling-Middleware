#ifndef PERSIST_H
#define PERSIST_H

#include <iostream>
#include <string>

using namespace std;

#ifdef __cplusplus
extern "C"
{
#endif

    void Create_File(int id);
    void Record(int id, const char *val, const char *opt_name, int up_value);
    bool undo(int r_id, const char *undo_update, int undo_val, int r_LT);
    const char *process_ret(std::string ret_val);
    bool stat_undo(int prev_updates);
    void Remote_Record(int id, int rid, const char *opt_name, int up_value);

#ifdef __cplusplus
}
#endif

#endif