(defun tevlis (m a) 
  (cond ((null m) '())
	('t (cons (eval (car m) a)
		  (tevlis (cdr m) a)))))
