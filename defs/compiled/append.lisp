(defun append (x y) 
  (cond ((null x) y)
	('t (cons (car x) (append (cdr x) y)))))
