#include <iostream>
#include <fstream>
#include <deque>
#include <vector>
using namespace std;

ifstream finput;

extern "C"
{
    string multiple_undo(const char *filename, int opt_nums)
    {

        finput.open(filename);
        vector<string> lines;
        string currentLine;

        // Read lines from the file and store the last n lines
        while (getline(finput, currentLine))
        {
            lines.push_back(currentLine);
        }
        finput.close();

        string last_opt_lines = lines[lines.size() - opt_nums - 1];
        return last_opt_lines;
    }
}