(defun next (t)
  ((lambda (max current)
     (cond
      ((eq max current) (err (list 'max max)))
      ('t (list max (inc current)))))
   (car t) (cadr t)))

