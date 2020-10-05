// gcc test.c -o binary
#include <stdio.h>
#include <unistd.h>

void init(){
    setvbuf(stdout, NULL, _IONBF, 0);
    setvbuf(stderr, NULL, _IONBF, 0);
    setvbuf(stdin, NULL, _IONBF, 0);
}


void checkTTY(){
    if (isatty(fileno(stdin)))
        printf( "stdin is a terminal\n" );
    else
        printf( "stdin is a file or a pipe\n");

    if (isatty(fileno(stdout)))
        printf( "stdout is a terminal\n" );
    else
        printf( "stdout is a file or a pipe\n");

    if (isatty(fileno(stderr)))
        printf( "stderr is a terminal\n" );
    else
        printf( "stderr is a file or a pipe\n");
}


int main(){
    // init();
    checkTTY();
    puts("Hello World");
    // sleep(5);
    puts("I woke up");
    return 0;
}