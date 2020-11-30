(defun mapcar (op args)
  (cond
   ((eq args '()) ())
   ('t (cons
	(funcall op (car args))
	(mapcar op (cdr args))))))
