 #include <stdio.h>

int sum(int a, int b) {
    int result;
    __asm__ ("addl %%ebx, %%eax;" : "=a" (result) : "a" (a), "b" (b));
    return result;
}

int main() {
    printf("%d\n", sum(2, 3));
    return 0;
}
