(defun factorial (n)
  (cond ((eq '0 n) '1)
	('t (mult n (factorial (minus n '1))))))
