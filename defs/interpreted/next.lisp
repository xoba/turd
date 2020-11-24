(defun next (t)
  (cond
   ((eq (car t) (cadr t)) (err (list 'max (car t))))
   ('t (list (car t) (inc (cadr t))))))

;; can't properly compile lambda's yet:
;;
;;(defun next (t)
;;  ((lambda (max current)
;;     (cond
;;      ((eq max current) (err (list 'max max)))
;;      ('t (list max (inc current)))))
;;   (car t) (cadr t)))
   
