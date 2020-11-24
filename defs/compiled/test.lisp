(defun test (x)
  ((lambda (first rest) 
     (list first rest))
   (car x) (cdr x)))

