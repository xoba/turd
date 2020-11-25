(defun test4 (x)
  ((label f (lambda (first rest)
	      (cond
	       ((eq first '0) (list first rest))
	       ('t (f (minus first '1) rest)))))
   (car x) (cdr x)))
