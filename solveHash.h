#ifndef SOLVEHASH_H
#define SOLVEHASH_H

#include <stdint.h>

char *solveHash(uint8_t *targetHash, char *startHashTrimmedLast, uint64_t startHashTwo, uint64_t radix16HexNumber, uint64_t shiftedNumber, uint64_t hh, uint64_t aa, uint64_t ff);

#endif /* SOLVEHASH_H */
