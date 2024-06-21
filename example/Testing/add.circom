pragma circom 2.0.0;

template test(){
    signal input one;
    var a[10];
    if(one == 1)
        a[1] = 1;
}

component main = test();