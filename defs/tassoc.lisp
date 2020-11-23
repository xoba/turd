(defun tassoc (t x y) 
  (cond
   ((eq (caar y) x) (cadar y))
   ('t (tassoc (next t) x (cdr y)))))

