(defun mapcar (op arglist)
  (cond
   ((eq arglist '()) ())
   ('t (cons
	(funcall op (car arglist))
	(mapcar op (cdr arglist))))))
