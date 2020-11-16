(defun assoc (x y) 
  (cond ((eq (caar y) x) (cadar y))
	('t (assoc x (cdr y)))))
					; TODO: make this fail if there's no assoc to be had
