((lambda y (f)
   ((lambda (x) (f (x x)))
    (lambda (x) (f (x x)))))
 (lambda (fact)
   (lambda (n)
     (cond
      ((eq '0 n) '1)
      ('t (mult n (fact (minus n '1))))))))
