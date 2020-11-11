((eval '((label firstatom (lambda (x)
			   (cond ((atom x) x)
				 ('t (firstatom (car x))))))
	y)
      '((y ((a b) (c d)))))
a
)
