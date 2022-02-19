import org.junit.Test;

import java.util.ArrayList;

import static org.junit.Assert.*;

public class OrderBookTest {

    @Test
    public void placeBuy() {
        ArrayList<Transaction> ts = new ArrayList<>();

        OrderBook b = new OrderBook("GOOG");
        b.placeBuy(1, 100, "David", 1);
        b.placeBuy(5, 100, "Brian", 2);
        b.placeBuy(20, 100, "Arjun", 3);
        b.placeBuy(4, 101, "Kevin", 4);

        b.placeSell(2, 101, "Andrew", 5);
        ts.add(new Transaction(
                "Andrew", "Kevin", "GOOG",
                2, 101, -1));

        b.placeSell(4, 100, "Bob", 6);
        ts.add(new Transaction(
                "Bob", "Kevin", "GOOG",
                2, 100, -1
        ));
        ts.add(new Transaction(
                "Bob", "David", "GOOG",
                1, 100, -1
        ));
        ts.add(new Transaction(
                "Bob", "Brian", "GOOG",
                1, 100, -1
        ));
        assertEquals(100, b.getMarketPrice(), .1);

        b.placeSell(5, 98, "Jake", 7);
        ts.add(new Transaction(
                "Jake", "Brian", "GOOG",
                4, 98, -1
        ));
        ts.add(new Transaction(
                "Jake", "Arjun", "GOOG",
                1, 98, -1
        ));
        assertEquals(98, b.getMarketPrice(), .1);

        ArrayList<Transaction> list = b.getTransactions();
        for (int i = 0; i < ts.size(); i++) {
            assertTrue(transEqual(list.get(i), ts.get(i)));
        }
    }

    @Test
    public void placeSell() {
        ArrayList<Transaction> ts = new ArrayList<>();

        OrderBook b = new OrderBook("GOOG");
        b.placeSell(1, 100, "David", 1);
        b.placeSell(5, 100, "Brian", 2);
        b.placeSell(20, 100, "Arjun", 3);
        b.placeSell(4, 99, "Kevin", 4);

        b.placeBuy(2, 99, "Andrew", 5);
        ts.add(new Transaction(
                "Kevin", "Andrew", "GOOG",
                2, 99, -1));

        b.placeBuy(4, 100, "Bob", 6);
        ts.add(new Transaction(
                "Kevin", "Bob", "GOOG",
                2, 100, -1
        ));
        ts.add(new Transaction(
                "David", "Bob", "GOOG",
                1, 100, -1
        ));
        ts.add(new Transaction(
                "Brian", "Bob", "GOOG",
                1, 100, -1
        ));
        assertEquals(100, b.getMarketPrice(), .1);

        b.placeBuy(5, 103, "Jake", 7);
        ts.add(new Transaction(
                "Brian", "Jake", "GOOG",
                4, 103, -1
        ));
        ts.add(new Transaction(
                "Arjun", "Jake", "GOOG",
                1, 103, -1
        ));
        assertEquals(103, b.getMarketPrice(), .1);

        ArrayList<Transaction> list = b.getTransactions();
        for (int i = 0; i < ts.size(); i++) {
            assertTrue(transEqual(list.get(i), ts.get(i)));
        }
    }

    @Test
    public void getMarketPrice() {}

    public boolean transEqual(Transaction a, Transaction b) {
        return a.security.equals(b.security) && a.amt == b.amt
                && a.price == b.price && a.from.equals(b.from)
                && a.to.equals(b.to);
    }
}