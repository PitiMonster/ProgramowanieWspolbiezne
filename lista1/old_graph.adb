with Ada.Text_IO;               use Ada.Text_IO;
with Ada.Numerics.Float_Random; 
with Ada.Containers.Indefinite_Vectors;
with Ada.Containers.Vectors;
with Ada.Integer_Text_IO;       use Ada.Integer_Text_IO;
with Ada.Numerics.Discrete_Random;

procedure Old_Graph is

    task type Observer is
        entry Observe (msg : String);
    end Observer;

    task type Node (Id : Integer) is
        entry InitGen;
        entry RunNode (Pkg : Integer);
    end Node;

    type Node_Type is access Node;
    
    

    package Node_Vector is new Ada.Containers.Indefinite_Vectors
       (Index_Type => Positive, Element_Type => Node_Type);

    package Integer_Vectors is
    new Ada.Containers.Vectors
    (Index_Type => Positive,
    Element_Type => Integer);

    package Nested_Vectors is
    new Ada.Containers.Vectors
    (Index_Type => Positive,
    Element_Type => Integer_Vectors.Vector,
    "=" => Integer_Vectors."=");

    NodesNum, PackagesNum, AddEdgesNum : Integer;
    O  :  Observer;
    Nodes : Node_Vector.Vector;
    Neighbors : Nested_Vectors.Vector;
    VisitedNodes : Nested_Vectors.Vector;
    HandledPacks : Nested_Vectors.Vector;
    ReceivedPacks : Integer;
    subtype R_Node is Integer range 1 .. 100000;
    package Rand_Node is new Ada.Numerics.Discrete_Random (R_Node);
    use Rand_Node; 
    G : Generator;
    temp : Integer;

    -- funkcja zwracająca losową liczbę z przedziału (Min, Max)
    function Get_Index(Min, Max : in Integer; Gen : in Generator) return Integer is 
    begin 
        temp := Integer(Random (Gen)) mod (Max-Min+1) + Min;
        return temp;
    end Get_Index; 



    task body Observer is
    begin
        loop
            select
                accept Observe (msg : String) do
                    -- wypisanie podusmowania
                    delay Duration (Float(0.4));

                    Put_Line (msg); -- wypisanie wiadomości przez Observera

                    if ReceivedPacks = PackagesNum then

                        Put_Line("");
                        Put_Line("");
                        Put_Line("Podsumowanie:");
                        Put_Line("");
                        Put_Line("Wierzchołki:");
                        Put_Line("");
                        for I in 1 .. Integer(Nodes.Length) loop
                            Put_Line("Paczki wierzchołka " & Integer'Image(I));
                            for J in 1 .. Integer(HandledPacks(I).Length) loop
                                Put(Integer'Image(HandledPacks(I)(J)) & ", ");
                            end loop;
                        Put_Line("");
                        Put_Line("");
                        end loop;

                        Put_Line("");
                        Put_Line("");
                        Put_Line("Paczki:");
                        Put_Line("");
                        for I in 1 .. PackagesNum loop
                            Put_Line("Wierzchołki paczki " & Integer'Image(I));
                            for J in 1 .. Integer(VisitedNodes(I).Length) loop
                                Put(Integer'Image(VisitedNodes(I)(J)) & ", ");
                            end loop;
                        Put_Line("");
                        Put_Line("");
                        end loop;

                    end if;
                end Observe;
            or
                terminate;
            end select;
        end loop;
    end Observer;
       
    task body Node is
        G   : Generator;
        Empty : Boolean;
        NeightId : Integer;
        CurrPackage : Integer;
    begin
        Empty := true; 
        loop
            select
                -- stowrzenie generatora dla node'a
                accept InitGen do
                    Reset (G, Id);
                    if Id /= NodesNum then 
                        Nodes(Id+1).InitGen; 
                    end if; 
                end InitGen;
            or
                -- 
                when Empty =>
                 accept RunNode (Pkg : Integer) do 
                    CurrPackage := Pkg;
                    Empty := false;
                end RunNode;
                if Empty = false then
                    O.Observe
                       ("Paczka numer" & Integer'Image (CurrPackage) &
                        " jest w wierzchołku " & Integer'Image (Id));
                    HandledPacks(Id).Append(CurrPackage);
                    VisitedNodes(CurrPackage).Append(Id);

                    delay Duration (Float( Get_Index(1000, 2000, G)/1000));
                    if Integer(Neighbors(Id).Length) /= 0 then
                        NeightId := Get_Index(1, Integer(Neighbors(Id).Length), G);
                        Nodes(Neighbors(Id)(NeightId)).RunNode(CurrPackage);
                    else 
                        ReceivedPacks := ReceivedPacks + 1;
                        delay Duration (Float( Get_Index(1000, 2000, G)/1000));
                        O.Observe("Paczka numer" & Integer'Image(CurrPackage) & " została odebrana");
                    end if;


                    Empty := true; 
                end if;

            or
                terminate;
            end select;
        end loop;
    end Node;

    NodeTemp, NeighTemp : Integer; 
    V : Integer_Vectors.Vector; 

begin
    Ada.Text_IO.Put ("Liczba wierzchołków: ");
    Ada.Integer_Text_IO.Get (NodesNum);
       Ada.Text_IO.Put ("Liczba losowych połączeń: ");
    Ada.Integer_Text_IO.Get (AddEdgesNum);
    Ada.Text_IO.Put ("Liczba paczek: ");
    Ada.Integer_Text_IO.Get (PackagesNum);

    begin
        -- inicjalizacja tablic
        for I in 1 .. NodesNum loop
            Nodes.Append (new Node (Id => I));
            Neighbors.Append(V);
            HandledPacks.Append(V);
        end loop;

        for I in 1 .. PackagesNum loop
            VisitedNodes.Append(V);
        end loop;

        for I in 1 .. NodesNum-1 loop
            Neighbors(I).Append(Integer(I+1));
        end loop;

        -- tworzenie powiązań między node'ami
        for I in 1 .. AddEdgesNum loop
            NodeTemp := Get_Index(1, NodesNum-2, G); 
            NeighTemp := Get_Index(NodeTemp+2, PackagesNum, G);
            while Neighbors(NodeTemp).Contains(NeighTemp) loop 
                NodeTemp := Get_Index(1, PackagesNum-2, G); 
                NeighTemp := Get_Index(NodeTemp+2, PackagesNum, G);
            end loop; 
            Neighbors(NodeTemp).Append(Integer(NeighTemp));
          --  Put_Line(Integer'Image(NodeTemp) & Integer'Image(NeighTemp));
        end loop; 

        Put_Line("Drukowanie grafu:");
        Put_Line("");
        for I in 1 .. Integer(Neighbors.Length) loop
            Put(Integer'Image(I) & " -> [");
            for J in 1 .. Integer(Neighbors(I).Length) loop
                Put(Integer'Image(Neighbors(I)(J)) & ", ");
            end loop;
            Put("]");
            Put_Line("");
        end loop;
        Put_Line("");


        Nodes(1).InitGen;
        for I in  1..PackagesNum loop
            Nodes(1).RunNode(I);
            delay Duration (Float( Get_Index(100, 200, G)/1000));
        end loop; 

    end;
end Old_Graph;

