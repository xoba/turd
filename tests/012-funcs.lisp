(((lambda (x y) (cond (((lambda (x) (eq x '())) x) y) ('t (cons (car x) (append (cdr x) y))))) '(c d))
(a b c d)
)
