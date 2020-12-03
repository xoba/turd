(defun tevcon (t c a) 
  (cond ((teval (next t) (caar c) a)
	 (teval (next t) (cadar c) a))
	('t (tevcon (next t) (cdr c) a))))
