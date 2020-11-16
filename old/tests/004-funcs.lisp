((
(label append
       (lambda (x y)
	 (cond ((
(label null
       (lambda (x) (eq x '())))
 x) y)
	       ('t (cons (car x) (append (cdr x)))))))
 '() '(c d))
(c d)
)
