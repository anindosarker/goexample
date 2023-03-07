#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <time.h>
#include <math.h>
#include "openssl/sha.h"

#define MaxRetries 10000000000

char* solveHash(uint8_t *targetHash, int64_t startHashTwo, int64_t startHashTrimmedLast, int64_t radix16HexNumber, int64_t shiftedNumber, int64_t hh, int64_t aa, int64_t ff) {
    uint64_t tries = 0, gg = 1;
    char *solvedHash = NULL;
    uint8_t notFound = 1;
    time_t ts;

    while(notFound) {
        if(!notFound)
            break;

        tries++;

        uint64_t c = hh * gg;

        for(ff; ff < startHashTwo; ff++) {
            c--;
            ts = time(NULL);
            tries++;

            if (tries > MaxRetries) {
                return "MaxRetries";
            }

            char *p1;
            asprintf(&p1, "%llx", radix16HexNumber+(ff>>(aa<<2)));
            char zeroString[17] = "0000000000000000";
            char *basep2;
            asprintf(&basep2, "%llx", ff&shiftedNumber);
            char *p2 = &basep2[strlen(basep2)-aa];

            char *g;
            asprintf(&g, "%lld%s%s", startHashTrimmedLast, p1, p2);

            uint8_t hash[32];
            SHA256_CTX ctx;
            SHA256_Init(&ctx);
            SHA256_Update(&ctx, g, strlen(g));
            SHA256_Final(hash, &ctx);

            uint8_t isEqual = 1;
            for (int i = 0; i < 32; i++) {
                if (hash[i] != targetHash[i]) {
                    isEqual = 0;
                    break;
                }
            }

            if (isEqual) {
                solvedHash = g;
                notFound = 0;
                break;
            }

            free(p1);
            free(basep2);
            free(g);
        }

        if (time(NULL) - ts <= 50) {
            gg++;
        } else {
            gg--;
            gg = fmax(gg, 1);
        }
    }

    return solvedHash;
}
