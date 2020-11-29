(defun factorial (n)
  (cond ((eq '0 n) '1)
	('t (mul n (factorial (sub n '1))))))
