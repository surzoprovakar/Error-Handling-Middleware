#include <iostream>
#include <fstream>
#include <algorithm>
using namespace std;

ifstream fin;

extern "C"
{
    // Need to make is string in set data structure
    string single_undo(const char *filename, int opt_nums)
    {

        fin.open(filename);
        string lastLine, currentLine, secondToLastLine;

        // Read lines from the file until the end
        while (getline(fin, currentLine))
        {
            secondToLastLine = lastLine;
            lastLine = currentLine;
        }
        fin.close();
        return secondToLastLine;
    }
}