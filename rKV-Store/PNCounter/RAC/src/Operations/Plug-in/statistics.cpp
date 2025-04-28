#include <iostream>
#include <cmath>
using namespace std;

extern "C"
{
    float calculateMean(int data[], int size)
    {
        float sum = 0.0, mean, standardDeviation = 0.0;

        for (int i = 0; i < size; i++)
        {
            sum += data[i];
        }

        return (sum / size);
    }

    float calculateSD(int data[], int size)
    {
        float sum = 0.0, mean, standardDeviation = 0.0;
        mean = calculateMean(data, size);

        for (int i = 0; i < size; i++)
        {
            standardDeviation += pow(data[i] - mean, 2);
        }

        return sqrt(standardDeviation / size);
    }

    bool should_undo(int *vals, int size)
    {
        // cout << "mean is " << calculateMean(vals, size) << endl;
        return calculateMean(vals, size) > 3 ? true : false;
    }
}