(defun tevlis (t m a) 
  (cond ((null m) '())
	('t (cons (teval (next t) (car m) a)
		  (tevlis (next t) (cdr m) a)))))
