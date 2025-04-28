#include <iostream>
#include <string.h>
#include <openssl/sha.h>
using namespace std;

extern "C"
{
    string calculateSHA256(int data[], int size)
    {
        SHA256_CTX sha256Context;
        SHA256_Init(&sha256Context);
        SHA256_Update(&sha256Context, &data, size);

        unsigned char sha256Hash[SHA256_DIGEST_LENGTH];
        SHA256_Final(sha256Hash, &sha256Context);

        string sha256Hex;
        for (int i = 0; i < SHA256_DIGEST_LENGTH; i++)
        {
            char hex[3];
            snprintf(hex, sizeof(hex), "%02x", sha256Hash[i]);
            sha256Hex += hex;
        }

        return sha256Hex;
    }

    bool should_undo(int *vals, int size)
    {
        return calculateSHA256(vals, size).length() != 64 ? true : false;
    }
}