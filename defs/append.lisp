(defun xappend (x y) 
  (cond ((null x) y)
	('t (cons (car x) (xappend (cdr x) y)))))
