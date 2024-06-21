pragma circom 2.0.0;

template AND() {
    signal input a;
    signal input b;
    signal output c;
    
    // Perform the AND operation using arithmetic constraints
    0 === a*(a-1);
    c <== a * b;
}

template multiAnd(n) {
    signal input in[n];
    signal output rv;
    component and;
    component ands[2];
    var i;
    
    if (n == 1) {
        rv <== in[0];
    } else if (n == 2) {
        and = AND();
        and.a <== in[0];
        and.b <== in[1];
        rv <== and.c;
    } else {
        and = AND();
        var n1 = n / 2;
        var n2 = n - n1;
        ands[0] = multiAnd(n1);
        ands[1] = multiAnd(n2);
        for (i = 0; i < n1; i++) ands[0].in[i] <== in[i];
        for (i = 0; i < n2; i++) ands[1].in[i] <== in[i + n1];
        and.a <== ands[0].rv;
        and.b <== ands[1].rv;
        rv <== and.c;
    }
}

component main = multiAnd(4);
