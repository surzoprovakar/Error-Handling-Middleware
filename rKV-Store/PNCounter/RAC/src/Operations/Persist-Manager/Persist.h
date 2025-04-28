#ifndef PERSIST_H
#define PERSIST_H

#include <iostream>
#include <string>

using namespace std;

#ifdef __cplusplus
extern "C"
{
#endif

    void Create_File(const char *id);
    void Record(const char *id, const char *val);
    const char *undo(const char *id, int opt_nums);
    const char *process_ret(string ret_val);
    bool stat_undo(int prev_updates);
    bool db_data_corrupt_undo(int prev_updates);

#ifdef __cplusplus
}
#endif

#endif