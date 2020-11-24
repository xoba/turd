(defun evcon (c a) 
  (cond ((eval (caar c) a)
	 (eval (cadar c) a))
	('t (evcon (cdr c) a))))
