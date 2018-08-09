import Data.List
import Data.List.Split (splitOn)
import System.Random
import Control.Monad (replicateM)

main = do
    seed <- newStdGen
    {-}
    (test_input, test_label) <- readCSV "mnist_test.csv"

    let input_dim = 784
    let output_dim = 10
    let x = take 8000 test_input 
    let y = take 8000 test_label 
    let x_val = drop 8000 test_input 
    let y_val = drop 8000 test_label 
    let model = multi_layer_net_initialize [input_dim,784,output_dim] seed 
    -}
    let input_dim = 784
    let output_dim = 10
    let x_dim = (200,input_dim)
    let y_dim = fst x_dim
    let x_val_dim = (500,input_dim)
    let y_val_dim =  fst x_val_dim
    let x = uncurry randomMatrix x_dim seed 
    let y = boundintlist (0,output_dim-1) $randomintlist y_dim seed 
    let x_val = uncurry randomMatrix x_val_dim seed 
    let y_val = boundintlist (0,output_dim-1) $ randomintlist y_val_dim seed
    let model = multi_layer_net_initialize [input_dim,784,output_dim] seed 
    train model x y x_val y_val  
    
readCSV csvpath = do
    contents <- readFile csvpath
    let csvdata = (map ( splitOn ",") (lines contents)) 
    let label = map readInt $ map head csvdata 
    let input = map (map readFloat ) $map tail csvdata 
   
    return (input, label)
readInt :: String -> Int 
readInt = read

readFloat :: String -> Float 
readFloat = read


affine_layer :: (Num a) => [[a]] -> [a] -> [[a]] -> ([[a]],([[a]],[[a]],[a])) 
affine_layer weight bias input = (out, (input, weight, bias))
    where
        out = madd (mmult input weight ) (unsqueeze res_dim bias)
        res_dim = (fst (dimension input), snd (dimension weight))

affine_layer_backward :: (Num a) => ([[a]],[[a]],[a])-> [[a]] -> ([[a]],[[a]],[a])  
affine_layer_backward (x,w,b) dout = (dx, dw, db)
    where
        dx = mmult dout (transpose w)
        dw = mmult (transpose x) dout
        db = squeeze(mmult (mones 1 (fst (dimension x))) dout)

relu_forward :: (Num a, Ord a) => [[a]] -> ([[a]], [[a]])
relu_forward x = (mapmat (max 0) x, x)
relu_backward :: (Num a, Ord a) => [[a]] -> [[a]] -> [[a]]
relu_backward x dout = mmult_elem (mapmat (\y -> if y <= 0 then 0 else 1) x) dout

softmax_loss :: ( Ord a,Floating a) => [[a]] -> [Int] -> (a,[[a]])
softmax_loss x y = (loss,dx)
    where
        probs = mapmat (exp) $madd x $maddinv $unsqueeze (dimension x) (map maximum x)
        probs_2 = mdiv_elem probs (unsqueeze (dimension probs) (reduce_sum probs 1))
        nN = fst (dimension x)
        loss = (-1.0 /(fromIntegral nN)) * sum (map log (zipWith (!!) probs_2 y))
        dx = smmult  (1.0 /(fromIntegral nN)) $madd probs_2 $maddinv $indexing_one (snd (dimension probs_2)) y 

multi_layer_net :: (Floating a, Ord a) => a ->[([[a]],[a])]->[[a]] ->[Int] ->  
    ((a,Double),[([[a]],[a])])
multi_layer_net reg model input y = (metric, grad)
    where
        (scores,cache) = mapAccumL (\(prevout) -> (\(weight,bias) -> 
            let (curout, newstate) = affine_layer weight bias prevout 
            in
                let (curout_relu, newstate_relu) = relu_forward curout
                in
                (curout_relu, (newstate_relu , newstate) )
                )) input model 
        (data_loss , dx )= softmax_loss scores y 
        loss = 0.5 * reg * (foldl (\reg_loss -> (\(weight, bias) -> 
            (reg_loss + sum (reduce_sum (mmult_elem weight weight) 1))
            )) 0 model) + data_loss
        maxindex xs = head $ elemIndices (maximum xs) xs
        accuracy = (\c -> (sum c)/(fromIntegral (length c))) $zipWith (\a -> \b -> if a == b then 1.0 else 0.0)(map maxindex scores) y
        metric = (loss, accuracy * 100) -- percent notation

        (_,grad,_) = foldr (\(weight, bias) -> (\(dprev_out , grad_cache, cache_left) -> 
            let (dout_relu) = relu_backward (fst (last cache_left)) dprev_out
            in
                let (dout, dw,db) = affine_layer_backward (snd (last cache_left)) dout_relu
                in
                (dout, (dw, db): grad_cache,(init cache_left))
             )) (dx, [],cache) model 

multi_layer_net_initialize :: ( Num b) => [Int] -> StdGen -> [([[Float]],[b])]
multi_layer_net_initialize layer_dimensions seed = model
    where
        model = pl2lp wl bl 
        wl = map  weight_initialization_method (pl2lp (init layer_dimensions) (tail layer_dimensions))
        bl = map bias_initialization_method (tail layer_dimensions)
        pl2lp l1 l2 = if length l1 == length l2 
            then zipWith (\x -> \y -> (x,y)) l1 l2
            else error "trying 2 pairize unequal length lists"
        weight_initialization_method = \x -> smmult (sqrt(2.0 / (fromIntegral (fst x) + fromIntegral (snd x)))) ((uncurry randomMatrix x) seed) -- Xavier Initialization
        bias_initialization_method = \x -> replicate x 0 

stochastic_gradient :: (Num a) =>a -> [([[a]],[a])] -> [([[a]],[a])]-> [([[a]],[a])]
stochastic_gradient learning_rate model grad = updated_model
    where 
        updated_model = zipWith (update) model grad
        update (weight, bias) (weight_grad, bias_grad) = ((madd weight (smmult (-1 * learning_rate) weight_grad)),(zipWith (-) bias (map ((*) learning_rate) bias_grad) )) 



train :: (Floating a, Show a,Ord a) => [([[a]],[a])] -> [[a]] -> [Int] -> [[a]] -> [Int] -> IO ([([[a]],[a])])
train model x y x_val y_val = do
    let batch_size = 200
    let reg = 0.0
    let learning_rate = 0.001
    let number_of_epoch = 10
    let number_of_data = fst $dimension x
    let number_of_val_data = fst $dimension x_val
    let iterations_per_epoch = number_of_data `div` batch_size
    seed <- newStdGen
    let train_step model_before = model_after 
            where
                chosen_indices = boundintlist (0,number_of_data-1) $randomintlist batch_size seed
                x_batch = map (x !! ) chosen_indices
                y_batch = map (y !! ) chosen_indices  
                (_, train_grad) = multi_layer_net reg model_before x_batch y_batch
                model_after = stochastic_gradient learning_rate model_before train_grad
    let iteration n model_before = if n <= 0 then model_before else model_after 
            where
                --loop
                model_intermediate = iteration (n-1) model_before
                model_after = train_step model_intermediate
    
    let epoch m model_before = if m <= 0 then return  model_before else do
                --loop
                model_intermediate <- epoch (m-1) model_before 
                let model_after = iteration iterations_per_epoch model_intermediate
                --loss print 
                let (train_metric,_) = multi_layer_net reg model_after x_train y_train
                        where
                            chosen_indices = boundintlist (0,number_of_data-1) $randomintlist number_of_val_data seed
                            x_train = map (x !! ) chosen_indices
                            y_train = map (y !! ) chosen_indices
                let (validation_metric, _) = multi_layer_net reg model_after x_val y_val
                print $"training metric : " ++ (show train_metric) ++ " validation metric : " ++ (show validation_metric)
                return model_after 
    print "************ Start training model **********"
    epoch number_of_epoch model
    

--normaldist random code from Data.Random.Normal Module 
boxMuller :: Floating a => a -> a -> (a,a)
boxMuller u1 u2 = (r * cos t, r * sin t) where r = sqrt (-2 * log u1)
                                               t = 2 * pi * u2

boxMullers :: Floating a => [a] -> [a]
boxMullers (u1:u2:us) = n1:n2:boxMullers us where (n1,n2) = boxMuller u1 u2
boxMullers _          = []

normal :: (RandomGen g, Random a, Floating a) => g -> (a,g)
normal g0 = (fst $ boxMuller u1 u2, g2)
  where
     (u1,g1) = randomR (0,1) g0
     (u2,g2) = randomR (0,1) g1

-- numpy-like helper function     

mmult :: (Num a) => [[a]] -> [[a]] -> [[a]]
mmult a b = if snd (dimension a) == fst (dimension b)
    then
        let b' = transpose b in  
            --map (\x -> map (\y -> sum (zipWith (*) x y)) b' ) a
            [ [ sum $ zipWith (*) ar bc | bc <- b' ] | ar <- a ]
    else error $"matrix format does not comply for mmult : " ++ show (dimension a) ++ " to " ++ show (dimension b)
madd :: (Num a) => [[a]] -> [[a]] -> [[a]]
madd a b = if dimension a == dimension b
    then zipWith(zipWith (+) ) a b
    else error $"matrix format does not comply for madd : " ++ show (dimension a) ++ " to " ++ show (dimension b)
mmult_elem :: (Num a) => [[a]] -> [[a]] -> [[a]]
mmult_elem a b
    | dimension a == dimension b = zipWith (zipWith (*) ) a b
    | otherwise = error $ "Unable to elementwise multiply two matrices : " ++ (show (dimension a))  ++ " to " ++ (show (dimension b))
mdiv_elem :: (Fractional a) => [[a]] -> [[a]] -> [[a]]
mdiv_elem a b
    | dimension a == dimension b = zipWith (zipWith (/) ) a b
    | otherwise = error $ "Unable to elementwise divide  two matrices : "++ (show (dimension a))  ++ " to " ++ (show (dimension b))
mzeros :: (Num a ) => Int -> Int -> [[a]]
mzeros x y = replicate x $ replicate y 0
maddinv :: (Num a) => [[a]] -> [[a]]
maddinv a = map (map negate) a
mapmat :: (Num a) => (a -> a) -> [[a]] -> [[a]]
mapmat f mat = map (map (f)) mat 
smmult :: (Num a) => a -> [[a]] -> [[a]]
smmult a b = mapmat ((*)a) b
mones :: (Num a) => Int -> Int -> [[a]]
mones x y = replicate x $ replicate y 1
indexing_one :: (Num a) => Int -> [Int] -> [[a]]
indexing_one max_index label = map (\x -> (zerolist x) ++ [1] ++ (zerolist (max_index - 1 - x))) label 
    where
        zerolist n = replicate n 0 
reduce_sum :: (Num a) => [[a]] -> Int -> [a]
reduce_sum a b 
    | null (head a) = []
    | b == 1 = map (sum) a
    | b == 0 = (sum (map head a)) : (reduce_sum (map tail a) b)
    | otherwise = error "reduce_sum axis must be 0 or 1" 
dimension :: (Num a) => [[a]] -> (Int,Int)
dimension p = (length p, length (head p))

unsqueeze :: (Num a) => (Int, Int) -> [a] -> [[a]]
unsqueeze (a,b) x 
    | length x == b = replicate a x
    | length x == a = transpose $ replicate b x 
    | otherwise = error $ "Unable to unsqueeze length "++(show (length x))  ++ " to "++ (show (a,b))
squeeze :: (Num a ) => [[a]] -> [a]
squeeze x 
    | fst (dimension x) == 1 =  head x
    | snd (dimension x) == 1 = head (transpose x)
    | otherwise = error $ "Unable to squeeze 2d array with dimensions like " ++ (show (dimension x))
splitEvery :: Int -> [a] -> [[a]]
splitEvery _ [] = []
splitEvery n xs = as : splitEvery n bs 
  where (as,bs) = splitAt n xs


-- transpose is on Data.List already
    
--print helper function
printm :: (Num a, Show a) => [[a]] -> IO ()
printm p = do
    let (x,y) = dimension p
    putStrLn $ (show x) ++ " x " ++ (show y) ++ "size weight matrix"
    --mapM_ print (p)
    putStrLn ""
layer_print (w,b) = do
        printm w 
        --print b 
        print $(show (dimension b)) ++ "length bias vector"
        putStrLn ""
model_print model_in = do
    mapM_ layer_print model_in 

-- random element matrix / list
randomMatrix :: Int -> Int -> StdGen -> [[Float]] --normaldist mean 0, variance 1
randomMatrix n m seed = (splitEvery m)$ randomlist (n*m) seed
randomlist ::  Int -> StdGen -> [Float]
randomlist n = take n . unfoldr (Just . normal)
boundintlist:: (Int,Int) -> [Int] -> [Int]
boundintlist (low, high) list_in = map ((flip rem) (high - low + 1).abs) list_in
boundedrandomList :: (Int, Int) -> Int -> IO([Int])
boundedrandomList (lower_bound, greater_bound) n = replicateM n $ randomRIO (lower_bound,greater_bound)       
randomintlist ::  Int -> StdGen -> [Int]
randomintlist n = take n . unfoldr (Just . random)

