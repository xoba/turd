(defun next (t)
  (cond
   ((eq (car t) (cadr t)) (err (list 'max (car t))))
   ('t (list (car t) (inc (cadr t))))))
