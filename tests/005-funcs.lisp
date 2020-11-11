((
(label append
       (lambda (x y)
	 (cond ((null x) y)
	       ('t (cons (car x) (append (cdr x)))))))
 '(a b) '(c d))
(a b c d)
)
