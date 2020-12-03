(defun length (x)
  (cond
   ((atom x) '0)
   ('t (add '1 (length (cdr x))))))


