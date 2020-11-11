((
(label append
       (lambda (x y)
	 (cond ((null x) y)
	       ('t (cons (car x) (append (cdr x)))))))
 '() '(c d))
(c d)
)
