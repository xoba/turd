;; test compilation of recursive labels
(defun test3 (x)
  ((label fx (lambda (first rest)
	      (cond
	       ((eq first '0) (list first rest))
	       ('t (fx (sub first '1) rest)))))
   (car x) (cdr x)))
