// gcc test.c -o binary
#include <stdio.h>
#include <unistd.h>

void init(){
    setvbuf(stdout, NULL, _IONBF, 0);
    setvbuf(stderr, NULL, _IONBF, 0);
    setvbuf(stdin, NULL, _IONBF, 0);
}

int main(){
    init();
    puts("Hello World");
    sleep(5);
    puts("I woke up");
    return 0;
}