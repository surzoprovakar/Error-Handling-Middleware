#ifndef SINGLE_UNDO_H
#define SINGLE_UNDO_H

#include <iostream>
#include <string>

using namespace std;

#ifdef __cplusplus
extern "C"
{
#endif
    const char *single_undo(const char *filename, int opt_nums);

#ifdef __cplusplus
}
#endif

#endif