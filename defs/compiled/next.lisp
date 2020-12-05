(defun next (t)
  ((lambda (max current)
     (cond
      ((eq max current) (err t))
      ('t (list max (inc current)))))
   (car t) (cadr t)))

