package main

func tester(*[]int) {

}

func main() {
    b := make([]int, 5)
    for i:=0; i<10; i++ {
        b[i] = 0;
    }
    tester(b)
}
