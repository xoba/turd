(defun length (x)
  (cond
   ((atom x) '0)
   ('t (plus '1 (length (cdr x))))))


