(s
 '(lambda (f x)
    (cond
     ((eq x '0) '1)
     ('t (mul x (f f (sub x '1))))))
 '7)
