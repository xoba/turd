(defun evlis (m a) 
  (cond ((null m) '())
	('t (cons (eval (car m) a)
		  (evlis (cdr m) a)))))
