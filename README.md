### **Hybrid Strassen's**

We implement a hybrid version of Strassen’s algorithm in Go. 

In the process, we
analytically derive an optimal threshold (i.e. crossover point) for switching from a pure implementation
of Strassen’s algorithm to the well-known standard matrix multiplication algorithm. 

We then empirically
search for an optimal practical crossover point by experi-
menting on various input sizes and thresholds. 

Finally, we use Strassen’s algorithm to model and count
triangle paths in a randomly-generated graph of 1024 vertices.

### **Reproducing our Results**

```go run strassen.go ```

```go run strassen.go 0 3 example.txt ```
