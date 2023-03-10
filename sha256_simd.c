#include <stdint.h>
#include <string.h>
#include "sha256-simd/sha256-simd.h"

void sha256_simd_wrapper(uint8_t *input, size_t input_len, uint8_t *output) {
    Sum256(input, input_len, output);
}
