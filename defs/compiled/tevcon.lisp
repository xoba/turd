(defun tevcon (c a) 
  (cond ((eval (caar c) a)
	 (eval (cadar c) a))
	('t (tevcon (cdr c) a))))
