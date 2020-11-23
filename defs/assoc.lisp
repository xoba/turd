(defun assoc (x y) 
  (cond
   ((eq (caar y) x) (cadar y))
   ('t (assoc x (cdr y)))))

