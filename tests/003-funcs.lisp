((
(label append
       (lambda (x y)
	 (cond ((
(label null
       (lambda (x) (eq x '())))
 x) y)
	       ('t (cons (car x) (append (cdr x)))))))
 '(a b) '(c d))
(a b c d)
)
