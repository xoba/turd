;; testing compilation of label
(defun test2 (x)
  ((label f
	  (lambda (first rest) 
	    (list first rest)))
   (car x) (cdr x)))

