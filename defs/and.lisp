(defun and (x y) 
  (cond (x (cond (y 't) ('t ())))
	('t '()))) ; TODO: this erroneously returns the string "()"!!
